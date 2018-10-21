import os
import shutil
import time
from pathlib import Path

rootpath = Path(Path(__file__).parent)
input_path = rootpath/Path("input.txt")
output_path = rootpath/Path("output.txt")
srcpath = rootpath/Path("source")
destpath = rootpath/Path("dest")
resultpath = rootpath/Path("result")

# copy, compile and run file 
for file in Path(srcpath).iterdir():
    filepath = destpath/Path(file.stem)
    Path.mkdir(filepath)
    shutil.copyfile(file, filepath/Path(file.name))
    shutil.copyfile(input_path, filepath/Path("input.txt"))
    shutil.copyfile(output_path, filepath/Path("output.txt"))
    os.popen("cd  " + str(filepath) + "/&&g++ -w " + str(file.name) + "&&a")
    time.sleep(2)
    shutil.copyfile(filepath/Path("output.txt"), resultpath/Path(file.stem + ".txt"))

# check the answer
answer = ["21", "15", "42", "98", "3", "error", "3", "1", "-1", "-12", "-2", "10", "3796",
"error", "error", "error", "error", "3796", "-120", "error", "error", "1", "0","0.1","2","2e+16"]

for file in Path(resultpath).iterdir():
    with file.open() as output:
        with open(resultpath/Path(file.stem + "_result.txt"),'w') as result_score:
            sum = 0
            for a in answer:
                b = output.readline()
                b = b.rstrip('\n')
                if a == b:
                    result_score.write("Right!\n")
                    sum = sum + 5
                else:
                    result_score.write(a + " " + b + "\n")
            result_score.write(str(sum))

