FROM mcr.microsoft.com/devcontainers/go:dev-1

RUN apt-get update && apt-get install -y \
    bash-completion \
    sqlite3

ENV USERNAME=vscode

RUN go install mvdan.cc/gofumpt@v0.7.0
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.23.0
RUN chown -R $USERNAME /go

# Bash history
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/command_history/.bash_history" \
    && mkdir /command_history \
    && touch /command_history/.bash_history \
    && chown -R $USERNAME /command_history \
    && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"
