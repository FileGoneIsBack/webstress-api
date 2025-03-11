#!/bin/bash
set -e
sudo -v
#############################
# Configuration 
SQL_DUMP_PATH="./assets/dump.sql"    
CONFIG_JSON="./assets/config/config.json"  
TOTAL_STEPS=5
PROGRESS_FILE="./assets/branding/vendors/install_progress.txt"
STATUS_FILE="./assets/branding/vendors/install_status.txt"
#############################

# --- Functions to update progress and status ---
update_progress() {
    local curr=$(cat "$PROGRESS_FILE")
    echo $(( curr + 1 )) > "$PROGRESS_FILE"
}

update_status() {
    echo "$1" > "$STATUS_FILE"
}

# --- Function: Display a dynamic progress bar ---
progress_bar() {
    local bar_width=50 

    while true; do
        clear  # Clear the terminal screen at the start of each loop iteration
        
        curr=$(cat "$PROGRESS_FILE")
        percentage=$(( curr * 100 / TOTAL_STEPS ))

        # Build progress bar
        filled=$(( bar_width * percentage / 100 ))
        empty=$(( bar_width - filled ))
        bar=$(printf "%0.s█" $(seq 1 $filled))
        bar+=$(printf "%0.s▒" $(seq 1 $empty))

        # Randomize positions of the stars
        # Generate random positions for 5-6 stars
        declare -A stars
        for i in {1..6}; do
            stars[$i,"x"]=$(( RANDOM % 35 + 1 ))  # Random horizontal position (1-35)
            stars[$i,"y"]=$(( RANDOM % 7 + 1 ))   # Random vertical position (1-7)
        done

        # Print stars in their new random positions
        for i in {1..6}; do
            tput cup "${stars[$i,"y"]}" "${stars[$i,"x"]}"
            echo "⋆"
        done

        # Print banner
        tput cup 2 0; echo "  _____          _ _ _       _     _    "  
        tput cup 3 0; echo " |_   _|_      _(_) (_) __ _| |__ | |_  "  
        tput cup 4 0; echo "   | | \ \ /\ / / | | |/ _  |  _ \| __| "  
        tput cup 5 0; echo "   | |  \ V  V /| | | | (_| | | | | |_  "  
        tput cup 6 0; echo "   |_|   \_/\_/ |_|_|_|\__, |_| |_|\__| "  
        tput cup 7 0; echo "                       |___/            "  

        # Print progress bar and current task
        tput cup 9 0; echo "========================================="  
        tput cup 10 0; echo "Installation progress: $percentage% [$bar]"
        tput cup 11 0; echo "Current task: $(cat "$STATUS_FILE")"

        # Exit if progress is complete
        if [ "$curr" -ge "$TOTAL_STEPS" ]; then
            break
        fi
        sleep .5
    done
}


# --- Start: Display Welcome ASCII Art ---
cat << "EOF"
⋆         ⋆           ⋆
⋆ _____  ⋆       _ _ _       _ ⋆   _    
 |_   _|_   ⋆  _(_) (_) __ _| |__ | |_  
   | | \ \ /\ / / | | |/ _` | '_ \| __| 
   | |  \ V  V /| | | | (_| | | | | |_  ⋆
   |_|   \_/\_/ |_|_|_|\__, |_| |_|\__| 
⋆        ⋆        ⋆     |___/  ⋆       ⋆  
=========================================
EOF
echo "Starting installation..."
sleep 2

# --- Ensure script is running on Ubuntu ---
if ! grep -qi ubuntu /etc/os-release; then
    echo "This script is intended for Ubuntu. Exiting..."
    exit 1
fi



# --- Prompt for SQL credentials and database name ---
while true; do
    read -rp "Do you need a SQL user? (yes/no): " CREATE_SQL_USER
    if [[ "$CREATE_SQL_USER" == "yes" || "$CREATE_SQL_USER" == "no" ]]; then
        break
    else
        echo "Invalid input. Please enter 'yes' or 'no'."
    fi
done

if [[ "$CREATE_SQL_USER" == "yes" ]]; then
    read -rp "Enter SQL username to create: " SQLUSER
    read -rsp "Enter password for $SQLUSER: " SQLPASS
    echo
    read -rp "Enter new database name: " DBNAME
else
    SQLUSER=""
    SQLPASS=""
    DBNAME=""
fi

# --- Initialize progress and status files ---
echo 0 > "$PROGRESS_FILE"
echo "Starting installation tasks..." > "$STATUS_FILE"

# --- Run installation tasks in a background subshell ---
(
    # Task 1: Install Latest Go
    if command -v go &>/dev/null; then
        update_status "Golang is already installed, skipping."
        sleep .5
        update_progress
    else
        update_status "Installing latest Golang..."
        curl -LO "https://go.dev/dl/go1.24.1.linux-amd64.tar.gz"
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "go1.24.1.linux-amd64.tar.gz"
        rm "go1.24.1.linux-amd64.tar.gz"
        if ! grep -q "/usr/local/go/bin" ~/.profile; then
            echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
        fi
        update_progress
    fi
        if ! command -v jq &> /dev/null; then
        update_status "jq not found, installing..."
        sudo apt update
        sudo apt install -y jq
        update_progress
    else
        update_status "jq is already installed."
        sleep .5
        update_progress
    fi
    # Task 2: Install MySQL Server
    update_status "Installing MySQL Server..."
    sleep 1
    sudo apt-get update
    sudo apt-get install -y mysql-server
    update_progress

# Task 3: Setup SQL user and database
if [[ "$CREATE_SQL_USER" == "yes" ]]; then
    update_status "Creating SQL user and database..."
   sudo mysql <<EOF
CREATE USER IF NOT EXISTS '$SQLUSER'@'localhost' IDENTIFIED BY '$SQLPASS';
GRANT ALL PRIVILEGES ON *.* TO '$SQLUSER'@'localhost' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CREATE DATABASE IF NOT EXISTS $DBNAME;
EOF
else
    update_status "Skipping SQL user and database creation..."
fi
update_progress

    # Task 4: Import SQL dump and update JSON (only if dump and config files exist)
    if [[ "$CREATE_SQL_USER" == "yes" && -f "$SQL_DUMP_PATH" && -f "$CONFIG_JSON" ]]; then
        update_status "Importing SQL dump"
        sudo mysql "$DBNAME" < "$SQL_DUMP_PATH"
        update_progress
        update_status "Updating JSON"
        jq --arg user "$SQLUSER" --arg pass "$SQLPASS" --arg table "$DBNAME" \
        '.username = $user | .password = $pass | .table = $table' \
        "$CONFIG_JSON" > tmp.$$.json && mv tmp.$$.json "$CONFIG_JSON"
    else
        update_status "Skipping SQL dump and JSON update..."
    fi
        update_progress
        sleep .5

    update_status "Installation complete!"
    update_progress
    sleep .5
) &

progress_bar
clear
echo "Starting website using 'go run .' if issuse please use 'CGO_ENABLED=1 go run .'"
echo "username, password, and db files have been updated in config.json!"
sleep 1
exec go run .
