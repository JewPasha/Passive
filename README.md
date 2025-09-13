# Passive Search Tool

Welcome to the Passive Search Tool! This tool allows you to search for user information based on full name, IP address, and usernames across various platforms such as Reddit, GitHub, and GitLab.

Make sure you have Go installed to use it.

### Running the Tool

To run the tool, use the `go run` command followed by the desired flags. Here are the available options:

- `-fn` : Search with full name
- `-ip` : Search with IP address
- `-u`  : Search with username
- `-help` : Display help message

### Examples

1. **Search by Full Name**:

    ```sh
    go run main.go -fn "Jean Dupont"
    ```

2. **Search by IP Address**:

    ```sh
    go run main.go -ip "8.8.8.8"
    ```

3. **Search by Username**:

    ```sh
    go run main.go -u "@example_username"
    ```

    The tool will search for the username across multiple platforms (Reddit, GitHub, GitLab, etc.) and display the results.

### Saving Results

Results are saved in the `results` directory. Each result is stored in a uniquely named text file within this directory.

## Features

- **Full Name Search**: Retrieves user information based on the full name from a local JSON database.
- **IP Address Search**: Provides information based on the provided IP address (e.g., ISP, city, country).
- **Username Search**: Checks the availability of a username across multiple platforms (Reddit, GitHub, GitLab).

**Full Name Search can be done only on dummy database because there are no public database with such personal information as it is illegal. 
