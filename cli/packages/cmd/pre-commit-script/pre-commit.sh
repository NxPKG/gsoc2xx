#!/bin/sh

# MANAGED BY GSOC2 CLI (Do not modify): START
gsoc2ScanEnabled=$(git config --bool hooks.gsoc2-scan)

if [ "$gsoc2ScanEnabled" != "false" ]; then
    gsoc2 scan git-changes -v --staged
    exitCode=$?
    if [ $exitCode -eq 1 ]; then
        echo "Commit blocked: Gsoc2 scan has uncovered secrets in your git commit"
        echo "To disable the Gsoc2 scan precommit hook run the following command:"
        echo ""
        echo "    git config hooks.gsoc2-scan false"
        echo ""
        exit 1
    fi
else
    echo 'Warning: gsoc2 scan precommit disabled'
fi
# MANAGED BY GSOC2 CLI (Do not modify): END