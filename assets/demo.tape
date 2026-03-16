Output assets/demo.gif
Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

# Init
Type "stl init"
Sleep 500ms
Enter
Sleep 2s

# Accept "Open browser?" confirm
Enter
Sleep 3s


# Select project: Down to "acme-api", then Enter
Down
Sleep 300ms
Enter
Sleep 2s

# Accept openapi spec path default
Enter
Sleep 1s

# Accept stainless config path default
Enter
Sleep 2s

# Accept target output paths (one Enter per target)
Enter
Sleep 500ms
Enter
Sleep 500ms
Enter
Sleep 8s

# Preview
Type "stl preview"
Sleep 500ms
Enter
Sleep 15s
Ctrl+C
Sleep 1s
