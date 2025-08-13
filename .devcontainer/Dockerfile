FROM mcr.microsoft.com/devcontainers/go:dev-1

# Install additional packages
RUN apt-get update && apt-get install -y \
    bash-completion \
    sqlite3

ENV USERNAME=vscode

# Install Go tools as the vscode user
USER $USERNAME
RUN go install mvdan.cc/gofumpt@v0.8.0 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 \
    && go install golang.org/x/tools/cmd/goimports@v0.36.0 \
    && go install github.com/onsi/ginkgo/v2/ginkgo@v2.23.4

# Switch back to root for final setup
USER root

# Bash history setup
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/command_history/.bash_history" \
    && mkdir /command_history \
    && touch /command_history/.bash_history \
    && chown -R $USERNAME:$USERNAME /command_history \
    && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"

# Set the default user
USER $USERNAME
