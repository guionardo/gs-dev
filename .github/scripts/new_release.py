import subprocess

from packaging.version import Version

print("# New release\n")


result = subprocess.run(["git", "describe", "--tags"], stdout=subprocess.PIPE)

last_tag = Version(result.stdout.decode().removeprefix("v").strip())

print(f"Current tag: v{last_tag}\n")
mj, mn, rl = last_tag.release

next_tag = Version(f"{mj}.{mn}.{rl + 1}")

print(f"Next tag:    v{next_tag}\n")

if not (nt := input(f"Inform new version or press enter to use {next_tag}:")):
    nt = str(next_tag)

next_tag = Version(nt)
if next_tag < last_tag:
    print("Next tag should be after last tag")
    exit(1)

print("\nRun this command to generate new release...\n")

cmd = f"git tag -a v{next_tag} -m 'v{next_tag}' && git push origin v{next_tag}"
print(cmd)

if input("Run this command now?").upper()[:1] in ["Y", "S"]:
    cmds = cmd.split("&&")
    for cmd in cmds:
        print(cmd)
        result = subprocess.run(cmd.split(), stdout=subprocess.PIPE)
        print(result.stdout.decode())
        if result.returncode:
            print(f"ERROR: {result.returncode}")
            print(result.stderr.decode())
            exit(1)
