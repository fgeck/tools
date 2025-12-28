# VHS demo for tools - A command-line bookmark manager
# https://github.com/charmbracelet/vhs

Output demos/demo.gif

Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 200ms

# Show the TUI
Type "./tools"
Enter
Sleep 1s

# Navigate through examples to show line wrapping
Type "j"
Sleep 300ms

# Add a new example
Type "a"
Sleep 500ms

# Fill in the form
Set TypingSpeed 50ms
Type "lsof -i :8080"
Sleep 300ms
Tab
Type "lsof"
Sleep 300ms
Tab
Type "show all processes running at port 8080"
Sleep 300ms
Enter
Sleep 1s

# Edit an example to demonstrate edit mode
Type "j"
Type "e"
Sleep 300ms

# Modify description - delete existing text first
Tab
Tab
Sleep 300ms
Left 13
Sleep 200ms
Backspace 7
Sleep 200ms
Type "listening"
Sleep 1s
Enter
Sleep 1s
Type "q"

# Clear screen
Hide
Type "clear"
Enter
Sleep 300ms
Show


# Add a new command under lsof tool using cli
Type `./tools add -c 'lsof -t -i :8080 | xargs kill -9' -n lsof -d 'kill process listening at port 8080'`
Enter
Sleep 1s

# Show updated list
Type "./tools"
Enter
Type "j"
Type "j"
Sleep 1s
Enter
Copy "lsof -t -i :8080 | xargs kill -9"
Sleep 1s
Paste
Sleep 5s
