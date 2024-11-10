import requests

URL = "https://forms.example.com/v1/entries/"
PROJECT_ID = "my_project_id"
AUTH_TOKEN = "my_auth_token"


def main():
    headers = {"accept": "application/json", "Authorization": f"Bearer {AUTH_TOKEN}"}

    response = requests.get(URL + PROJECT_ID, headers=headers, timeout=15)

    if response.status_code != 200:
        print(
            f"Request failed with status code: {response.status_code}::{response.reason}"
        )
        return

    data = response.json()
    print(data)


if __name__ == "__main__":
    main()
