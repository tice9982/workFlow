## Download/Install

git clone the repository to `$GOPATH/src/workFlow`.

## Report Issues / Send Patches

## test
"""
        === RUN   TestSingleFlow
        workflow_test.go 14>>>  s on run
        workflow.go 93>>>       s cancel
        --- PASS: TestSingleFlow (2.00s)

        === RUN   TestSerialFlow
        workflow.go 142>>>      serial flow! children count:3
        workflow_test.go 36>>>  serial s 0 on enter
        workflow.go 96>>>       s 0 done
        workflow_test.go 42>>>  serial s 0 on run
        workflow.go 96>>>       s 0 done
        workflow_test.go 48>>>  serial s 0 on quit
        workflow.go 96>>>       s 0 done
        workflow_test.go 36>>>  serial s 1 on enter
        workflow.go 96>>>       s 1 done
        workflow_test.go 42>>>  serial s 1 on run
        workflow.go 93>>>       s 1 cancel
        --- PASS: TestSerialFlow (4.00s)

        === RUN   TestParallelFlow

        workflow.go 203>>>      parallel flow! children count:3
        workflow_test.go 69>>>  parallel p 1 on enter
        workflow_test.go 69>>>  parallel p 2 on enter
        workflow_test.go 69>>>  parallel p 0 on enter
        workflow.go 96>>>       p 1 done
        workflow.go 96>>>       p 2 done
        workflow.go 96>>>       p 0 done
        workflow_test.go 74>>>  parallel p 1 on run
        workflow_test.go 74>>>  parallel p 0 on run
        workflow_test.go 74>>>  parallel p 2 on run
        workflow.go 96>>>       p 0 done
        workflow.go 96>>>       p 1 done
        workflow.go 96>>>       p 2 done
        workflow_test.go 81>>>  parallel p 0 on quit
        workflow_test.go 81>>>  parallel p 1 on quit
        workflow_test.go 81>>>  parallel p 2 on quit
        workflow.go 96>>>       p 0 done
        workflow.go 96>>>       p 1 done
        workflow.go 96>>>       p 2 done
        --- PASS: TestParallelFlow (3.00s)

        PASS
"""