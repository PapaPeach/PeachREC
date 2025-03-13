# PeachREC
PeachREC is an automated demo recorder for Team Fortress 2 tournament matches. It is a spiritual successor to PREC that makes clever use of animations and scripts to detect and record tournament matches.  
This differs from Valve's automated demo recording by NOT recording Casual matches, and ONLY recording scrims, matches, or PUGs.

# Installation
I recommend making a backup of your *autoexec.cfg* and custom HUD. This isn't strictly necessary, the installer should at most add one line to both of those, but put as much trust in my coding as feel comfortable with.
- Download the latest release from the [Releases page](https://github.com/PapaPeach/PeachREC/releases).
- Put **peachrec_installer.exe** in your *tf/custom* folder.
- Run **peachrec_installer.exe** and answer any prompts it gives you.
- If you did not allow the program to edit your *autoexec.cfg* you'll need to manually add `exec peachrec` to it or add `+exec peachrec` to your launch options.
- If you'd like to configure custom start/stop messages or sounds, you can do that in the generated *_PeachREC/cfg/peachrec.cfg* file.

That's the entire installation process. Though for added security I recommend using `bind [key] ds_status` in conjunction with `ds_notify 2` to display on your demo recording status on your screen when pressing the bound key.

# Notes
- PeachREC relies on being present **when the match starts** to record a demo automatically. If you join or have to reconnect halfway through a match **you must record a demo manually**.
- There is a known quirk with PeachREC where in a tournament server, joining a team, joining spectator, then joining a team again, will trigger a recording to start and immediately stop. It's inconsequential and probably impossible to fix.

# Compatibility
| **Category** | **Compatibility** |
| :--- | :--- |
| Windows | Fully compatible |
| Linux | Theoretically compatible |
| Mac | Theoretically compatible |
| Custom HUDs | Fully compatible |
| Default HUD | Fully compatible |
| Custom configs | Fully compatible |
| mastercomfig | Fully compatible |
| No config | Fully compatible |
| Casual servers | Fully compatible |
| Community servers | Fully compatible |
| Tournament servers | Fully compatible |
| Valve Competitive servers | Untested |
| Mann Vs Machine servers | Will automatically ready once |

# How's it work?
- PeachREC Checks if the current server allows the tournament setup menu to open (where you set team names and ready), if the menu opens the server is determined to be a match/scrim server.
- Then the script will be armed and will record whenever the player spawns AND is no longer able to open the tournament setup menu, because the match has started.
- PeachREC also checks if the player joins a new server to prevent falsely triggering a recording when going from a tournament server to a Casual server.

There's more complexitiy to the script, but that is to work around the quirkiness of HUD animations, match starts, respawns, etc.
