{{ template "header.tpl" . }}

<main>
    <h2>Please login</h2>
    <form action="/web/login" method="POST">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required>
        <br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required>
        <br>
        <button type="submit">Login</button>
    </form>
</main>

{{ template "footer.tpl" . }}