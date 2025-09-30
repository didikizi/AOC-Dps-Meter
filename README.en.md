# AOC DPS Meter

A Damage meter for Ashes of Creation game. This tool analyzes combat logs to provide detailed statistics about your performance in battles.

![AOC DPS Meter](screenshots/main-interface.png)

## ⚠️ Important Disclaimers

- **This project is an independent community initiative and is NOT affiliated with Intrepid Studios or Ashes of Creation**
- **The tool only reads and analyzes game log files; it does NOT modify the game client or interact with game servers**
- **It is intended for personal and educational use only**
- **Users are solely responsible for ensuring that their use of this tool complies with the game’s Terms of Service and End User License Agreement**
- **Real-time log monitoring may be subject to restrictions under the game’s policies; we recommend using this tool primarily for post-combat log analysis**
- **The developers provide this project “as is” and assume no responsibility for any consequences resulting from its use**

## 🚀 Features

- **Real-time DPS/HPS tracking** - Monitor your damage and healing per second
- **Combat statistics** - Detailed breakdown of hits, crits, and damage
- **Ability analysis** - See which abilities deal the most damage
- **Target tracking** - Analyze damage dealt to different enemies
- **Sortable tables** - Sort data by any column (damage, hits, crit rate, etc.)
- **Collapsible sections** - Hide/show abilities and targets tables
- **Precise damage numbers** - Exact damage values without rounding
- **Auto log detection** - Automatically finds AOC log files

## 📋 Requirements

- Windows 10/11

## 🎯 Roadmap (50 Stars Goal)

If this project reaches 50 stars, I plan to add:

1. **Group Interface** - Multi-player statistics tracking
2. **Loot Checker** - Analyze loot drops and rewards

## 📥 Installation

### Build from Source

#### Prerequisites
- Go 1.24.5 or later
- Node.js (for frontend)
- Wails v2
- Git

#### Build Steps

1. Clone the repository:
```bash
git clone https://github.com/didikizi/AOCDpsMetr.git
cd AOCDpsMetr
```

2. Install Wails:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

3. Build the application:
```bash
wails build
```

The executable will be created in the `build/bin/` directory.

#### Development Mode
To run the application in development mode with hot reload:
```bash
wails dev
```

#### Troubleshooting Build Issues
- Make sure Go is properly installed and in your PATH
- Ensure Wails is installed correctly: `wails doctor`
- Check that Node.js is installed for frontend dependencies
- On Windows, you may need Visual Studio Build Tools

## 🎮 Usage

1. **Launch the application**
2. **Start Game** and begin playing
3. **Click "Start Monitoring"** in the DPS meter
4. **Fight enemies** - the meter will track your combat data in real-time
5. **View statistics** - DPS, damage, crit rates, and more
6. **Sort tables** - Click column headers to sort data
7. **Collapse sections** - Use the ▼ buttons to hide/show tables

## 📊 Interface Overview

### Main Statistics
- **Max DPS** - Highest damage per second achieved
- **Current DPS** - Real-time damage per second
- **Total Damage** - Exact damage dealt (no rounding)
- **Hits/Crits** - Number of hits and critical hits
- **Crit Rate** - Percentage of critical hits

### Abilities Table
- Damage and healing per ability
- Hit counts and crit rates
- Sortable by any column

### Targets Table
- Damage dealt to each enemy
- Kill counts per target
- Detailed combat statistics

## 🔧 Troubleshooting

### FAQ

**Q: The application can't find AOC log files**

A: The application searches for logs in these locations:
- `%USERPROFILE%\AppData\Local\AOC\Saved\Logs\AOC.log`
- `%USERPROFILE%\Documents\AOC\Logs\AOC.log`
- `%USERPROFILE%\AppData\Roaming\AOC\Logs\AOC.log`

If your logs are in a different location, you can:
1. Modify the source code to add your custom path

**Q: Statistics show old data when starting monitoring**

A: The application processes all existing events from the log file when you start monitoring, so you see cumulative statistics from the beginning of your gaming session.

**Q: The application doesn't update in real-time**

A: Make sure:
- Game is running
- You're in combat (the meter only tracks during fights)
- The log file is being written to

**Q: Can I use this with other players?**

A: Currently, this tool only tracks your personal combat data. Group functionality is planned for the 50-star milestone.

## 🛠️ Development

### Project Structure
```
AOCDpsMetr/
├── frontend/          # Web interface (HTML/CSS/JS)
├── internal/          # Go backend
│   ├── app/          # Main application logic
│   ├── metrics/      # Statistics calculation
│   ├── parser/       # Log file parsing
│   └── watcher/      # File monitoring
├── build/            # Build output
└── main.go           # Application entry point
```

### Key Components
- **Parser** - Extracts combat events from AOC log files
- **Calculator** - Processes events and calculates statistics
- **Watcher** - Monitors log files for new events
- **Frontend** - Real-time UI with sorting and collapsing

## 📝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Wails v2](https://wails.io/) - Go + Web frontend framework
- Thanks to the community for feedback and testing

---

**Remember**: This tool is for personal use only. Use responsibly and in accordance with the game's terms of service.
