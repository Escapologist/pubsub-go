# GoPubSub

A simple social media platform built with Go that allows users to post messages, follow other users and interact with posts. This project is for learning and teaching Go web development concepts including HTTP servers, templates, and basic social networking features.

> **Note**: This project is loosely inspired by Raymond Hettinger's excellent [pubsub educational demo](https://github.com/rhettinger/modernpython/tree/master/pubsub) from his Modern Python course. While the core concepts are similar, this implementation is built in Go and includes additional features.

## Features

- User authentication
- Post creation and viewing
- User following system
- Profile pages
- Search functionality
- Hashtag support
- Responsive Bootstrap UI

## Prerequisites

- Go 1.23.0 or later
- A modern web browser

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go_pubsub.git
cd go_pubsub
```

2. Run the application:
```bash
go run main.go
```

3. Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

### Authentication
- Login with one of the demo accounts, i.e.:
  - Username: `rich`
  - Password: any non-empty string

### Creating Posts
1. Log in to your account
2. Use the text area at the top of the page to write your message
3. Click the "Post" button to publish
4. Use hashtags (e.g., #hello) to tag your posts

### Following Users
1. Visit a user's profile page by clicking their username
2. Click the "Follow" button to follow them
3. Their posts will appear in your feed

### Searching
1. Use the search bar in the navigation menu
2. Search for text in posts or hashtags
3. Click on hashtags in posts to see all posts with that tag

## Project Structure

```
go_pubsub/
├── main.go              # Application entry point
├── web/                 # Web handlers and routing
│   └── web.go
├── pubsub/             # Core business logic
│   └── pubsub.go
├── session/            # Session management
│   └── session.go
├── templates/          # HTML templates
│   ├── base.html
│   ├── main.html
│   ├── user.html
│   └── ...
└── static/            # Static assets
    └── ...
```

## Development

The project uses Go modules for dependency management. The main dependencies are:
- Standard library `net/http` for web server
- Bootstrap 4.3.1 for UI components
- jQuery 3.3.1 for JavaScript functionality

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
