# bmj-gha-poc
```
name: Tester-Run
on: [push, pull_request]
jobs:
  build:
    name: Get Last Version Number
    runs-on: <your runner>
    steps:
      - uses: actions/checkout@v1

      - name: BMJ Custom Github Action
        uses: BMJ-Ltd/bmj-gha-poc@1.0.13
        id: version
        with:
          ecr_name: <ecr_name>
      - name: Get the version
        run: |
          echo "${{ steps.version.outputs.myOutput }}"
```          