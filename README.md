
# BMJ

This GHA will check a given repository and increment a version number

The version number must be in the format of major.minot.patch like 1.2.3

If no such version is found on any of the images that are already in the ECR a new version number 0.0.1 
will be created by default.




## Environment Variables

To run this project, you will need to have all the correct GitHub secrets set to access the ECR in AWS


## Usage/Examples

The two required parameters are the ***ecr_name*** and the ***patch***

Valid parameters for patch are 
* major
* minor
* patch

Where ***major*** will increase the ***major*** version and set the ***minor*** and ***patch*** to zero

Where ***minor*** will increase the ***minor*** version and set the ***patch*** to zero

Where ***patch*** will increase the ***patch*** version

### Add this to your gitHubWorkflow


```

       ## Version Number increase
      - name: BMJ Version number increment
        uses: BMJ-Ltd/xxxxx@1.0.40
        id: version
        with:
          ecr_name: ${{ github.event.repository.name }}
          version_type: "patch"
      - name: Get the version
        run: |
          echo '${{ steps.version.outputs.newVersion }}'

```


# Hi, I'm Mike! 👋


## 🚀 About Me
I'm a DevOps Engineer in BMJ
