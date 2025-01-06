# Mirror Purge

## Description
The Mirror Purge is a utility tool designed to identify duplicate files on your system. It helps free up storage space and organize your files more effectively by locating and flagging identical files.

## Features
- Scans directories and subdirectories for duplicate files.
- Identifies duplicates using:
  - File names
  - File sizes
  - Hash comparisons (e.g., MD5, SHA-256).
- Outputs results in an easy-to-read format.
- Offers options to delete or move duplicates.

## Installation
1. Clone this repository:
   ```bash
   git clone https://github.com/IsahiRea/mirror-purge.git
   ```
2. Navigate to the project directory:
   ```bash
   cd mirror-purge
   ```
3. Build the project:
   ```bash
   go build
   ```

## Usage
1. Run the tool:
   ```bash
   ./mirror-purge [options] <directory>
   ```
2. Options:
   - `--hash <hash-type>` or `-h <hash-type>`: Use hash comparisons (md5, sha256).
   - `--output <file>` or `-o <file>`: Specify an output file for the results.
   - `--delete` or `-d`: Prompt to delete duplicates.

Example:
```bash
./mirror-purge -h sha256 -o results.txt ~/Documents
```

## Roadmap
- Add GUI for non-technical users.
- Support for additional file types.
- Advanced filtering options (e.g., by file type or date).