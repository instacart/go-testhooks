#!/usr/bin/env python3
import os
import sys
import re

TESTHOOK_SOURCE_FILE_RE = re.compile(r"testhooks/[^/.]+\.go")


def tmp_dir():
    return sys.argv[1]


def get_content_lines_numbered(fname):
    with open(fname) as f:
        return list(enumerate(f.read().splitlines()))


def get_run_section(lines):
    start = [i for i, line in lines if line.startswith("main.run ")][0]
    end = [i for i, line in lines if i > start and not line.startswith("\t")][0]
    return lines[start:end]


def run_check(fname, func):
    lines = get_content_lines_numbered(os.path.join(tmp_dir(), fname))
    run_section = get_run_section(lines)
    error_message = func(run_section)

    if error_message:
        print(
            "ERROR: %s\n\nmain.run section is shown below:\n%s\n\nFull output from go build:\n%s\n"
            % (
                error_message,
                "\n".join(["%04d: %s" % (i, line) for i, line in run_section]),
                "\n".join(["%04d: %s" % (i, line) for i, line in lines]),
            ),
            file=sys.stderr,
        )
        sys.exit(1)


def run_checks(bintest):
    def check_default(run_section):
        if not any(
            [line for i, line in run_section if TESTHOOK_SOURCE_FILE_RE.search(line)]
        ):
            return (
                "disassemblychecks/%s: main.run built without notest tag should have references to testhooks/*.go"
                % bintest
            )

    run_check(
        bintest + "/go-build-logs-default.txt",
        check_default,
    )

    def check_notest(run_section):
        line_numbers = [
            i for i, line in run_section if TESTHOOK_SOURCE_FILE_RE.search(line)
        ]
        if any(line_numbers):
            return (
                "disassemblychecks/%s: main.run built with notest tag should have no references to testhooks/*.go, but it has on line %d"
                % (bintest, line_numbers[0])
            )

    run_check(
        bintest + "/go-build-logs-notest.txt",
        check_notest,
    )


def main():
    run_checks("errhook")
    run_checks("valuehook")


if __name__ == "__main__":
    main()
