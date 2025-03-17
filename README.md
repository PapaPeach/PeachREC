# PeachREC
Here's a video walkthrough and demo of the PeachREC installer:  
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
- I **HIGHLY** recommend using `bind [key] ds_status` in conjunction with `ds_notify 2` to display on your demo recording status on your screen when pressing the bound key.

That's the entire installation process. If you'd like to test PeachREC after installing, instructions are below.

### Important Notes:
- PeachREC relies on being present **when the match starts** to record a demo automatically. If you join or have to reconnect halfway through a match **you MUST record a demo manually**.
- If you install a different HUD after running **peachrec_installer.exe**, you will need to re-run **peachrec_installer.exe** so that it can modify your new HUD accordingly.
- PeachREC makes use of Valve's improved demo recording utilizing `ds_...` commands. Therefore, it retains custom prefixes allowing for users to organize demos however they find most intuitive. You can also use it along side map-prefix generators such as [this one](https://www.teamfortress.tv/47180/demo-support-ds-prefix-on-any-map).

# Known Issues
| **ID** | **Suspected Cause Description** | **Type** | **Fix** | **Resolution Status** |
| :---: | --- | --- | --- | --- |
| 1 | Join server within WaitingForPlayers time (first 30s). | False negative | Hit resupply bind.<br>Self-resolves next respawn. | League / server host needs to add mp_waitingforplayers_cancel 1 to server config.<br>Discussing with league admins / server hosts. |
| 2 | Match starts while in AFK status. | False negative | Hit resupply bind.<br>Self-resolves next respawn. | Investigating issue. |
| 3 | Server re-executes the config. | False positive | Hit resupply bind.<br>Self-resolves next respawn. | Likely unfixable. |
| 4 | In a tournament server, join a team, join spectator, and join a team again. | False positive | Self-resolves instantly | Likely unfixable and mostly inconsequential. |
| 5 | Lag spike or unlucky timing on match start. | Escape key is disabled | Open and close your text chat. | Testing alternative timings on PeachRecClose animation. |
| 6 | Join MvM server and resupply / change class before readying. | Auto-ready once | Manually unready.<br>After readying PeachREC will identify the MvM session. | Likely unfixable and mostly inconsequential. |
- **False Negative** - Does not begin recording a demo when one should have been. **Top Priority!**
- **False Positive** - Begins recording a demo at an incorrect time. **Mild inconvenience.**

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

# Testing
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
