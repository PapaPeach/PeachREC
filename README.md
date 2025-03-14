# PeachREC
[![Watch the video](https://img.youtube.com/vi/GuldUj6zqXs/hqdefault.jpg)](https://youtu.be/GuldUj6zqXs)  
PeachREC is an automated demo recorder for Team Fortress 2 tournament matches. It is a spiritual successor to PREC that makes clever use of animations and scripts to detect and record tournament matches.  
This differs from Valve's automated demo recording by NOT recording Casual matches, and ONLY recording scrims, matches, or PUGs. It will also automatically start recording at the beginning of the actual match (not pregame) and stop recording at the end of the actual match.

# Installation
I recommend making a backup of your *autoexec.cfg* and custom HUD. This isn't strictly necessary, the installer should at most add one line to both of those, but put as much trust in my coding as feel comfortable with.
- Download the latest release from the [Releases page](https://github.com/PapaPeach/PeachREC/releases).
- Put **peachrec_installer.exe** in your *tf/custom* folder.
- Run **peachrec_installer.exe** and answer any prompts it gives you.
- If you did not allow the program to edit your *autoexec.cfg* you'll need to manually add `exec peachrec` to it or add `+exec peachrec` to your launch options.
- If you'd like to configure custom start/stop messages or sounds, you can do that in the generated *_PeachREC/cfg/peachrec.cfg* file.

That's the entire installation process. Though for added security I recommend using `bind [key] ds_status` in conjunction with `ds_notify 2` to display on your demo recording status on your screen when pressing the bound key.

### Testing
The easiest and most reliable way to test PeachREC is to play a demo file that includes the pregame and halftime of a tournament match. This won't trigger PeachREC to record a demo but it will trigger start and stop sounds and console feedback.  
Offline testing is a bit wonky because of TF2's quirkiness. But to test the functionality of PeachREC on an offline / private server:
- Start on a server with `mp_tournamnet 0` (this is usually default).
- Set `mp_tournament 1` and make sure the team ready statuses appear at the top of the screen, if they don't you may have to set `mp_waitingforplayers_restart 1`. Once those are visible you should hear the PeachREC standby chime and see "=====PeachREC.waiting.for.match.to.start=====" in console.
- Set `mp_tournament 0` this simulates the match having started and no longer being able to open the tournament setup menu.  
**PeachREC will not record yet.**
- Change class **once** or use a resupply bind to respawn **twice**. In a match this would simulate being sent back to the spawn rooms when pregame ends. For whatever reason resupply binds are only half as wonky as changing class or a match starting.  
**PeachREC will now start recording.**
- Set `mp_tournament 1` and make sure the team ready statuses appear at the top of the screen, if they don't you may have to set `mp_waitingforplayers_restart 1`. Once those are visible PeachREC will stop recording.
- You can repeat this process if you'd like, you may need to resupply bind or respawn via other means inbetween each full run of the test to simulate the halftime.

# Notes
- PeachREC relies on being present **when the match starts** to record a demo automatically. If you join or have to reconnect halfway through a match **you must record a demo manually**.
- If you install a different HUD after running **peachrec_installer.exe**, you will need to re-run **peachrec_installer.exe** so that it can modify your new HUD accordingly.
- PeachREC makes use of Valve's improved demo recording utilizing `ds_...` commands. Therefore, it retains custom prefixes allowing for users to organize demos however they find most intuitive. You can also use it along side map-prefix generators such as [this one](https://www.teamfortress.tv/47180/demo-support-ds-prefix-on-any-map).
- There is a known quirk where in MvM respawning for the first time in a match will trigger you to automatically ready once. But after that initial ready PeachREC will detect that you're in an MvM server and will behave itself.
- There is a known quirk with PeachREC where in a tournament server, joining a team, joining spectator, then joining a team again, will trigger a recording to start and immediately stop. It's inconsequential and probably impossible to fix.

# Compatibility
| **Category** | **Compatibility** |
| :--- | :--- |
| Windows | Fully compatible |
| Linux | Fully compatible |
| Mac | Theoretically compatible |
| Custom HUDs | Fully compatible |
| Default HUD | Fully compatible |
| Custom configs | Fully compatible |
| mastercomfig (latest) | Fully compatible |
| No config | Fully compatible |
| Casual servers | Fully compatible |
| Community servers | Fully compatible |
| Tournament servers | Fully compatible |
| Valve Competitive servers | Untested |
| Mann Vs Machine servers | Will automatically ready once |

# How Does It Work?
- The installer searches your custom HUD for hudanimations_manifest.txt, and searches the files referenced within it for animations pertinent to PeachREC.
- If/when the installer finds a relevant animation, it makes a copy of that animation's code so that PeachREC's animations will behave identically to your custom HUD's, with PeachREC code added in addition to your HUD's.
- The installer then inserts a file path to the PeachREC animations at the top spot in your HUD's hudanimations_manifest.txt file, so your hud will load PeachREC.
- Then the installer generates peachrec.cfg with the necessary code and appends `exec peachrec` to your autoexec.cfg if you allow it to.
- If you don't have a custom HUD or a custom config, the installer will generate PeachREC to be self sufficient using default HUD's values.

In game the script works as follows:
- PeachREC Checks if the current server allows the tournament setup menu to open (where you set team names and ready), if the menu opens the server is determined to be a match/scrim server.
- Then the script will be armed and will record whenever the player spawns AND is no longer able to open the tournament setup menu, because the match has started.
- PeachREC also checks if the player joins a new server to prevent falsely triggering a recording when going from a tournament server to a Casual server.

There's more complexity to the script, but that is to work around the quirkiness of HUD animations, match starts, respawns, etc.

# Is This a Virus?
Nope. Your anti-virus software may flag it as a potential virus, likely just because it is an executable file that your anti-virus software has never seen before. Your anti-virus is doing its job by alerting you to this, running untrusted executable files is inherently dangerous. You should review the code and determine my trustworthiness for yourself and decide whether to use my installer or not.

If you don't trust me you can still use the installer by reviewing the code yourself in this GitHub repository. It is written in Go so it'll be pretty clear to anyone with programming experience what is going on.  
Once you're satisfied with your code review you can build the installer from source using any updated version of Go. All libraries used are native Go libraries so you can build it by opening a terminal window in the directory you have the source code in and executing a `go build` command, which will build you your own **peachrec_installer.exe** to use following my original install instructions.
