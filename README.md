# SnipForge - AI-Powered Code Snippet Generator

![GitHub](https://img.shields.io/github/license/peetya/snipforge-cli)

> **SnipForge** is a powerful command-line interface (CLI) tool that utilizes OpenAI's GPT technology to generate code 
> snippets for various programming and tooling languages based on a given description. It aims to save time and effort 
> for developers by providing a convenient way to generate code on-demand.


## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Binary Installation](#binary-installation)
  - [Local Installation](#local-installation)
- [Usage](#usage)
- [Flags](#flags)
- [Example](#example)
  - [Basic example](#basic-example)
  - [Advanced example](#advanced-example)
- [Important Note on Generated Output](#important-note-on-generated-output)
- [Contributing](#contributing)
- [License](#license)

## Features

- Leverages OpenAI's GPT technology for intelligent code snippet generation
- Supports multiple programming and tooling languages
- Matches your specific needs with customizable options
- Guides you through the code generation process in interactive mode step-by-step 
- Outputs code snippets to `stdout` or saves them to a file

## Installation

### Binary Installation

1. Download the appropriate binary for your operating system from the [latest release](https://github.com/peetya/snipforge-cli/releases/latest) page.
2. Extract the downloaded archive and place the `snipforge` binary in a directory that is included in your system's `PATH`.
3. Open a terminal and type `snipforge version` to confirm the installation was successful.

### Local Installation

If you prefer to run SnipForge from source, follow these steps:

1. Clone the repository:
```bash
$ git clone https://github.com/peetya/snipforge-cli.git
```

2. Change into the `snipforge-cli` directory:
```bash
$ cd snipforge-cli
```

3. Install the necessary dependencies:
```bash
$ go mod download
```

4. Run:
```bash
$ go run main.go version
```

## Usage

To get started with **SnipForge**, install the CLI tool and run the `generate` command, providing the required flags and 
options. Alternatively, you can start the interactive mode by running the `generate` command without any flags, which 
will guide you through the code generation process step-by-step.

```bash
$ snipforge generate [flags]
```

For more information on available commands and flags, refer to the help output by running:

```bash
$ snipforge generate --help
```

## Flags

Here's a detailed explanation of the available flags for the `generate` command:

```bash
-g, --goal:             The functionality description for the code snippet
-l, --language:         The programming or tooling language to generate code in (e.g., PHP, Golang, etc.)
-v, --language-version: The version of the programming or tooling language to generate code for (if applicable)
-o, --output:           The output file path for the generated code snippet
--output-chmod:         The chmod value to apply to the output file (default 644)
--stdout:               Print the generated code snippet to stdout instead of saving to a file
-k, --openai-key:       The OpenAI API key
-m, --openai-model:     The OpenAI model to use (default "gpt-3.5-turbo")
-q, --quiet:            Suppress all output except for the generated code snippet
-d, --dry-run:          Do not generate a code snippet, only print the generated description
```

## Example

### Basic example

Here's a basic example of how to use SnipForge to generate a Python code snippet for sorting a list of numbers:

```bash
$ snipforge generate --language python --language-version 3.11 --goal "sort a list of numbers" --output sorted_numbers.py
```

This command will generate a Python code snippet in the `sorted_numbers.py` file, with the goal of sorting a list of 
numbers.

```python
numbers = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5]
numbers.sort()
print(numbers)
```

### Advanced example

Here's an advanced example demonstrating how to use SnipForge to generate a PHP code snippet that fulfills more complex 
goals. In this example, we'll use SnipForge in interactive mode.

```bash
$ snipforge generate
```

First, we need to define a set of goals that will be used to generate the snippet.

```
Goal #1: A controller that returns a list of users via the "/api/v1/users" endpoint
Goal #2: The output format can be changed via content negotiation
Goal #3: Support pagination using the page and limit query parameters
Goal #4: Read the users from the injected UserRepositoryInterface
Goal #5: The controller must follow the PSR-12 coding standard
Goal #6: The controller must follow the PSR-4 autoloading standard
```

Next, we need to define the programming language and version to generate the snippet for.

```
Language: Symfony
LanguageVersion (optional): 6
```

Then we need to define the output path:

```
Output file path: src/Controller/Api/V1/UserController.php
```

Then it will generate the following code snippet for you in `src/Controller/Api/V1/UserController.php`:

```php
<?php

// src/Controller/Api/V1/UserController.php

namespace App\Controller\Api\V1;

use App\Repository\UserRepositoryInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/api/v1/users')]
class UserController extends AbstractController
{
    public function __construct(private UserRepositoryInterface $userRepository)
    {
    }

    #[Route('', name: 'api_v1_users_index', methods: ['GET'])]
    public function index(Request $request): Response
    {
        $page = $request->query->getInt('page', 1);
        $limit = $request->query->getInt('limit', 10);

        $users = $this->userRepository->getPaginatedUsers($page, $limit);

        $serializer = $this->container->get('serializer');
        $content = $serializer->serialize(['data' => $users], $request->getPreferredFormat(), ['groups' => ['user']]

        return new Response($content, 200, ['Content-Type' => $request->getMimeType($request->getPreferredFormat())]);
    }
}
```

## Important Note on Generated Output

While SnipForge does its best to generate accurate and functional code snippets, it's crucial to review the generated 
output before using it in your projects. AI-generated code may sometimes contain errors or inconsistencies, so always 
double-check the results to ensure correctness.

## Contributing

SnipForge is an open-source project, and we welcome contributions from developers like you! Feel free to submit issues, 
suggest new features, or create pull requests to help improve the project. We appreciate your support and collaboration. 
:heart:

## License

SnipForge is released under the [MIT License](LICENSE).
