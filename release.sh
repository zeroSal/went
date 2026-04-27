#!/bin/bash
set -e

CYAN="\033[0;36m"
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
RESET="\033[0m"

current_branch=$(git rev-parse --abbrev-ref HEAD)

echo -e "[*] Checking current branch..."
if [[ "$current_branch" != "staging" ]]; then
  echo -e "\n${RED}[x] You must be on 'staging' to run this script. Current branch: '$current_branch'.${RESET}"
  exit 1
fi

echo -e "[*] Checking the current version..."
latest_tag=$(git describe --tags --abbrev=0 || echo "0.0.0")
echo -e "${CYAN}[i] Latest tag: $latest_tag${RESET}"

read -p "[?] New version (e.g., 1.2.3): " new_version
if [[ -z "$new_version" ]]; then
  echo -e "${RED}[x] No version entered.${RESET}"
  exit 1
fi

if [[ ! "$new_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo -e "${RED}[x] Version must match X.Y.Z pattern.${RESET}"
  exit 1
fi

echo -e "\n[*] Checking for existing tag..."
if git rev-parse "$new_version" >/dev/null 2>&1 || git ls-remote --exit-code origin "refs/tags/$new_version" >/dev/null 2>&1; then
  echo -e "${RED}[x] Tag '$new_version' already exists.${RESET}"
  exit 1
fi

echo -e "\n[*] Tagging..."
git tag -a "$new_version" -m "Version $new_version"
git push origin "$new_version"

echo -e "\n${GREEN}[v] Release '$new_version' created successfully.${RESET}"
echo -e "\n${CYAN}[i] Now merge 'staging' into 'main' and push:${RESET}"
cat <<EOF
  git checkout main
  git merge staging
  git push
EOF
