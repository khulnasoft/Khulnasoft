from khulnasoft import Sandbox

# Timeout 3s for the sandbox to open
sandbox = Sandbox(template="base", timeout=3)  # $HighlightLine

sandbox.close()
