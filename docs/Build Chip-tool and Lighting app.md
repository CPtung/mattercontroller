# Build Chip-tool and Lighting app

# 環境

- x64_86 Debian 11

## 前置作業

- 安裝相依套件
    - 參考 Matter github building guide
    
    [connectedhomeip/docs/guides/BUILDING.md at master · project-chip/connectedhomeip](https://github.com/project-chip/connectedhomeip/blob/master/docs/guides/BUILDING.md)
    
- 下載 Matter Github repository
    - git clone [git@github.com](mailto:git@github.com):project-chip/connectedhomeip.git
- 建立環境
    - cd connectedhomeip
    - `git submodule update --init`
    - sh ./scripts/bootstrap.sh
        - screenshot
            
            ```bash
            ➜  connectedhomeip git:(master) sh ./scripts/bootstrap.sh
            2025-09-17 06:34:27,692 Loading extra packages for darwin
            2025-09-17 06:34:27,692 Appending: /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json for this platform
            2025-09-17 06:34:27,692 Skipping: windows (i.e. /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json)
            -e
              WELCOME TO...
            
                     █
                     █
                 ▄   █   ▄                                █     █
                 ▀▀█████▀▀      ▄▀▀▀▄ ▄▀▀▀▄    ▄▀▀▀▀▄█  ▀▀█▀▀▀▀▀█▀▀   ▄▀▀▀▀▄    ▄▀▀
               ▀█▄       ▄█▀   █     █     █  █      █    █     █    █▄▄▄▄▄▄█  █
                 ▀█▄   ▄█▀     █     █     █  █      █    █     █    █         █
              ▄██▀▀█   █▀▀██▄  █     █     █   ▀▄▄▄▄▀█    ▀▄▄   ▀▄▄   ▀▄▄▄▄▀   █
             ▀▀    █   █    ▀▀
            
            -e   BOOTSTRAP! Bootstrap may take a few minutes; please be patient.
            
            Downloading and installing packages into local source directory:
            
              Setting up CIPD package manager...done (1m45.4s)
              Setting up Project actions........skipped (0.1s)
              Setting up Python environment.....done (2m38.8s)
              Setting up pw packages............skipped (0.1s)
              Setting up Host tools.............done (0.1s)
            
            Activating environment (setting environment variables):
            
            -n   Setting environment variables for CIPD package manager...
            done
            -n   Setting environment variables for Project actions........
            skipped
            -n   Setting environment variables for Python environment.....
            done
            -n   Setting environment variables for pw packages............
            skipped
            -n   Setting environment variables for Host tools.............
            done
            
            Checking the environment:
            
            ================================================================================
            The Pigweed developer tool (`pw`) uses Google Analytics to report usage,
            diagnostic, and error data. This data is used to help improve Pigweed, its
            libraries, and its tools.
            
            Telemetry is not sent on the very first run. To disable reporting of telemetry
            for future invocations, run this terminal command:
            
                pw cli-analytics --opt-out
            
            If you opt out of telemetry, no further information will be sent. This data is
            collected in accordance with the Google Privacy Policy
            (https://policies.google.com/privacy). For more details on how Pigweed collects
            telemetry, see https://pigweed.dev/pw_cli_analytics.
            ================================================================================
            20250917 06:39:02 INF enabling analytics
            20250917 06:39:02 INF Analytics enabled, run 'pw cli-analytics --opt-out' to disable.
            20250917 06:39:02 INF See https://pigweed.dev/pw_cli_analytics/ for more details.
            20250917 06:39:02 INF
            20250917 06:39:02 INF Per-project settings:
            20250917 06:39:02 INF     ~/Documents/Works/ubiquiti/connectedhomeip/pigweed.json
            20250917 06:39:02 INF
            20250917 06:39:02 INF Per-user per-project settings:
            20250917 06:39:02 INF     ~/Documents/Works/ubiquiti/connectedhomeip/.pw_cli_analytics.user.json
            20250917 06:39:02 INF
            20250917 06:39:02 INF Per-user settings:
            20250917 06:39:02 INF     ~/.pw_cli_analytics.json
            20250917 06:39:02 INF
            20250917 06:39:02 INF api_secret = m7q0D-9ETtKrGqHAcQK2kQ
            20250917 06:39:02 INF measurement_id = G-NY45VS0X1F
            20250917 06:39:02 INF debug_url = https://www.google-analytics.com/debug/mp/collect
            20250917 06:39:02 INF prod_url = https://www.google-analytics.com/mp/collect
            20250917 06:39:02 INF report_command_line = False
            20250917 06:39:02 INF report_project_name = False
            20250917 06:39:02 INF report_remote_url = False
            20250917 06:39:02 INF report_subcommand_name = limited
            20250917 06:39:02 INF uuid = e8f3398a-a5d5-4dc0-b93c-bfba1cede253
            20250917 06:39:02 INF enabled = False
            20250917 06:39:02 INF Environment passes all checks!
            ================================================================================
            The Pigweed developer tool (`pw`) uses Google Analytics to report usage,
            diagnostic, and error data. This data is used to help improve Pigweed, its
            libraries, and its tools.
            
            Telemetry is not sent on the very first run. To disable reporting of telemetry
            for future invocations, run this terminal command:
            
                pw cli-analytics --opt-out
            
            If you opt out of telemetry, no further information will be sent. This data is
            collected in accordance with the Google Privacy Policy
            (https://policies.google.com/privacy). For more details on how Pigweed collects
            telemetry, see https://pigweed.dev/pw_cli_analytics.
            ================================================================================
            
            Environment looks good, you are ready to go!
            
            To reactivate this environment in the future, run this in your
            terminal:
            
            -e   source ./activate.sh
            
            To deactivate this environment, run this:
            
            -e   deactivate
            
            Installing pip requirements for all...
            
            [notice] A new release of pip is available: 23.2.1 -> 25.2
            [notice] To update, run: pip install --upgrade pip
            /Users/justincptung/connectedhomeip/scripts/helpers/bash-completion.sh: line 52: syntax error near unexpected token `<'
            /Users/justincptung/connectedhomeip/scripts/helpers/bash-completion.sh: line 52: `                readarray -t COMPREPLY < <(compgen -W "$("$1" targets --format=completion "$cur")" -- "$cur")'
            ```
            
    - source scripts/activate.sh
        - screenshot
            
            ```bash
            ➜  connectedhomeip git:(master) source scripts/activate.sh
            2025-09-17 06:40:47,683 Loading extra packages for linux
            2025-09-17 06:40:47,683 Appending: /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json for this platform
            2025-09-17 06:40:47,683 Skipping: windows (i.e. /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json)
            
              WELCOME TO...
            
                     █
                     █
                 ▄   █   ▄                                █     █
                 ▀▀█████▀▀      ▄▀▀▀▄ ▄▀▀▀▄    ▄▀▀▀▀▄█  ▀▀█▀▀▀▀▀█▀▀   ▄▀▀▀▀▄    ▄▀▀
               ▀█▄       ▄█▀   █     █     █  █      █    █     █    █▄▄▄▄▄▄█  █
                 ▀█▄   ▄█▀     █     █     █  █      █    █     █    █         █
              ▄██▀▀█   █▀▀██▄  █     █     █   ▀▄▄▄▄▀█    ▀▄▄   ▀▄▄   ▀▄▄▄▄▀   █
             ▀▀    █   █    ▀▀
            
              ACTIVATOR! This sets your shell environment variables.
            
            Activating environment (setting environment variables):
            
              Setting environment variables for CIPD package manager...done
              Setting environment variables for Project actions........skipped
              Setting environment variables for Python environment.....done
              Setting environment variables for pw packages............skipped
              Setting environment variables for Host tools.............done
            
            Checking the environment:
            
            20250917 06:40:48 INF Environment passes all checks!
            
            Environment looks good, you are ready to go!
            ```
            

## 編譯 lighting-app

- chmod +x scripts
- ./scripts/examples/gn_build_example.sh examples/lighting-app/linux out/lighting-app

```bash
➜  connectedhomeip git:(master) chmod +x scripts
➜  connectedhomeip git:(master) ./scripts/examples/gn_build_example.sh examples/lighting-app/linux out/lighting-app
2025-09-17 06:42:03,915 Loading extra packages for darwin
2025-09-17 06:42:03,915 Appending: /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json for this platform
2025-09-17 06:42:03,916 Skipping: windows (i.e. /Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/pw_env_setup/py/pw_env_setup/cipd_setup/python311.json)

  WELCOME TO...

         █
         █
     ▄   █   ▄                                █     █
     ▀▀█████▀▀      ▄▀▀▀▄ ▄▀▀▀▄    ▄▀▀▀▀▄█  ▀▀█▀▀▀▀▀█▀▀   ▄▀▀▀▀▄    ▄▀▀
   ▀█▄       ▄█▀   █     █     █  █      █    █     █    █▄▄▄▄▄▄█  █
     ▀█▄   ▄█▀     █     █     █  █      █    █     █    █         █
  ▄██▀▀█   █▀▀██▄  █     █     █   ▀▄▄▄▄▀█    ▀▄▄   ▀▄▄   ▀▄▄▄▄▀   █
 ▀▀    █   █    ▀▀

  ACTIVATOR! This sets your shell environment variables.

Activating environment (setting environment variables):

  Setting environment variables for CIPD package manager...done
  Setting environment variables for Project actions........skipped
  Setting environment variables for Python environment.....done
  Setting environment variables for pw packages............skipped
  Setting environment variables for Host tools.............done

Checking the environment:

20250917 06:42:04 INF Environment passes all checks!

Environment looks good, you are ready to go!

+ env
PW_PIGWEED_CIPD_INSTALL_DIR=/Users/justincptung/connectedhomeip/.environment/cipd/packages/pigweed
PW_ZAP_CIPD_INSTALL_DIR=/Users/justincptung/connectedhomeip/.environment/cipd/packages/zap
CIPD_CACHE_DIR=/Users/justincptung/.cipd-cache-dir
TERM_PROGRAM=iTerm.app
PW_PACKAGE_ROOT=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/packages
TERM=xterm-256color
SHELL=/bin/zsh
TMPDIR=/var/folders/t7/lt9q70297qv3bc97kcv0c20c0000gn/T/
CONDA_SHLVL=0
TERM_PROGRAM_VERSION=3.4.10
OLDPWD=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip
TERM_SESSION_ID=w0t0p0:5CCC833D-5F97-498C-93D5-489B82B820AA
ZSH=/Users/justincptung/.oh-my-zsh
USER=justincptung
COMMAND_MODE=unix2003
PW_BRANDING_BANNER=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/scripts/setup/banner.txt
CONDA_EXE=/opt/anaconda3/bin/conda
SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.YrRBdBswbz/Listeners
__CF_USER_TEXT_ENCODING=0x0:2:53
_PW_ENVIRONMENT_CONFIG_FILE=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/scripts/setup/environment.json
VIRTUAL_ENV=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/pigweed-venv
PAGER=less
_CE_CONDA=
LSCOLORS=Gxfxcxdxbxegedabagacad
PATH=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo/out/host/host_tools:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/pigweed-venv/bin:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/arm/bin:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/arm:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/zap:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/pigweed/bin:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/pigweed:/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd:/opt/anaconda3/condabin:/opt/homebrew/bin:/Users/justincptung/.cargo/env:/Users/justincptung/go/bin:/Users/justincptung/bin:/usr/local/bin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/Users/justincptung/.cargo/bin:/Users/justincptung/.local/bin:/Users/justincptung/.lmstudio/bin:/Users/justincptung/.local/bin
PW_ROOT=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/third_party/pigweed/repo
PW_ARM_CIPD_INSTALL_DIR=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment/cipd/packages/arm
__CFBundleIdentifier=com.googlecode.iterm2
PWD=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip
LANG=zh_TW.UTF-8
ITERM_PROFILE=JustinStyle
XPC_FLAGS=0x0
_PW_ACTUAL_ENVIRONMENT_ROOT=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip/.environment
_CE_M=
XPC_SERVICE_NAME=0
SHLVL=2
HOME=/Users/justincptung
COLORFGBG=15;0
_PW_ROSETTA=0
LC_TERMINAL_VERSION=3.4.10
ITERM_SESSION_ID=w0t0p0:5CCC833D-5F97-498C-93D5-489B82B820AA
CONDA_PYTHON_EXE=/opt/anaconda3/bin/python
LESS=-R
LOGNAME=justincptung
PW_PROJECT_ROOT=/Users/justincptung/Documents/Works/ubiquiti/connectedhomeip
LC_TERMINAL=iTerm2
COLORTERM=truecolor
_=/usr/bin/env
+ gn gen --check --fail-on-unused-args --root=examples/lighting-app/linux out/lighting-app --args=
Done. Made 5778 targets from 510 files in 1260ms
+ ninja -C out/lighting-app
ninja: Entering directory `out/lighting-app'
[715/715] ld ./chip-lighting-app
```