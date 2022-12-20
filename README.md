# go-cookbook
Go-cookbook is a web application for hosting recipes for cooking.
## Functionalities
Functionalities in a quick overview:
- Store, adjust and delete recipes.
- Store conversions for ingredients, so ingredients are displayed in required measurements (cups, grams, milliliters, teaspoons, etc.).
- Adjust portions of existing recipes to the portions you need when cooking.
- Add tags to organize your recipes.

## How to run
- Clone/Download this repository.
- In command-line, go to the folder where you've cloned/unpacked the repository.
- Run `go test`, this will create a default login: admin/test (Note: in future want to add login in the main folder to create default user when none exists).
- Run `go build` to create executable.
- Run the executable (in Linux: `./go-cookbook`, in windows `go-cookbook.exe`).
- Open the cookbook in your browser at `localhost:8081`.

## More information
- Since this is a basis application, no server is used. All data is stored into json files, located in the config folder.
