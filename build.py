import subprocess
import sys


def main(argv):
  

    if(len (argv) ==1):
        arg=argv[0]

        if arg== "build":
            print("build")

        if arg== "run":
            command="go run main.go"
            process= subprocess.Popen(command.split(),stdout=sys.stdout)
            output,error = process.communicate()
            print(output)

    else:
        print(" build \t build the project")
        print(" run \t run the project locally")

main(sys.argv[1:])