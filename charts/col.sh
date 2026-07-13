#!/usr/bin/env bash



# Default directories/files to exclude

EXCLUDED_DIRS=("node_modules" ".git" "__pycache__" ".venv" "venv" "build" "dist" ".mypy_cache" ".env", "kuberun-agent

")

TARGET_DIR="."

OUTPUT_FILE=""



# Helper function to check if an item exists in the exclusion list

is_excluded() {

  local item="$1"

  for exc in "${EXCLUDED_DIRS[@]}"; do

    if [[ "$item" == "$exc" ]]; then

      return 0

    fi

  done

  return 1

}



# Parse command line arguments

while [[ $# -gt 0 ]]; do

  case "$1" in

    -o|--output)

      OUTPUT_FILE="$2"

      shift 2

      ;;

    -e|--exclude)

      # Split by comma in case of "-e .gitignore,noshow"

      IFS=',' read -ra SPLIT_ITEMS <<< "$2"

      for item in "${SPLIT_ITEMS[@]}"; do

        # Strip leading/trailing whitespace

        trimmed=$(echo "$item" | xargs)

        EXCLUDED_DIRS+=("$trimmed")

      done

      shift 2

      ;;

    -*)

      echo "Unknown option: $1" >&2

      exit 1

      ;;

    *)

      TARGET_DIR="$1"

      shift

      ;;

  esac

done



# Ensure the target directory exists

if [ ! -d "$TARGET_DIR" ]; then

  echo "Error: Directory '$TARGET_DIR' does not exist." >&2

  exit 1

fi



# Resolve absolute paths

TARGET_DIR=$(realpath "$TARGET_DIR")

if [ -z "$OUTPUT_FILE" ]; then

  OUTPUT_FILE="$TARGET_DIR/combined_files_output.txt"

else

  OUTPUT_FILE=$(realpath -m "$OUTPUT_FILE")

fi



# Banner message (Matching the Python rich styling)

echo -e "\nCollecting \033[36m$TARGET_DIR\033[0m → \033[1m$OUTPUT_FILE\033[0m\n"



# Clear or initialize the output file safely

: > "$OUTPUT_FILE"

files_written=0



# Recursively find all files, sorted alphabetically (null-terminated to handle spaces safely)

while IFS= read -r -d '' file; do

  abs_file=$(realpath "$file")



  # 1. Skip if it is the output file itself

  if [[ "$abs_file" == "$OUTPUT_FILE" ]]; then

    continue

  fi



  # 2. Skip if any part of the path is in the EXCLUDED_DIRS list

  skip_file=false

  IFS='/' read -ra PATH_PARTS <<< "$abs_file"

  for part in "${PATH_PARTS[@]}"; do

    if is_excluded "$part"; then

      skip_file=true

      break

    fi

  done



  if $skip_file; then

    continue

  fi



  # Get the path relative to the target directory

  rel_path="${abs_file#$TARGET_DIR/}"



  # Write the 80-character delimiter headers

  printf '%.0s=' {1..80} >> "$OUTPUT_FILE"

  echo -e "\nFILE PATH: $rel_path" >> "$OUTPUT_FILE"

  printf '%.0s=' {1..80} >> "$OUTPUT_FILE"

  echo "" >> "$OUTPUT_FILE"



  # 3. Read content safely (Check readability and binary status)

  if [ ! -r "$abs_file" ]; then

    echo "[binary or unreadable file — skipped]" >> "$OUTPUT_FILE"

  elif grep -qI . "$abs_file" 2>/dev/null || [ ! -s "$abs_file" ]; then

    # It's a valid text file or an empty readable file

    cat "$abs_file" >> "$OUTPUT_FILE" 2>/dev/null

  else

    # Failed text-check (likely binary)

    echo "[binary or unreadable file — skipped]" >> "$OUTPUT_FILE"

  fi



  # Add trailing newlines matching Python's out.write("\n\n")

  echo -e "\n\n" >> "$OUTPUT_FILE"

  ((files_written++))



done < <(find "$TARGET_DIR" -type f -print0 | sort -z)



# Success footer message

echo -e "\033[32mpassed\033[0m  $files_written files written to \033[1m$OUTPUT_FILE\033[0m"

