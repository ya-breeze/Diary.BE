{{ template "header.tpl" . }}

<main>
    {{ with .UserID }}
    {{ else }}
        <h2>This is the Home Page. Please login</h2>
        <form action="/" method="POST">
            <label for="username">Username:</label>
            <input type="text" id="username" name="username" required>
            <br>
            <label for="password">Password:</label>
            <input type="password" id="password" name="password" required>
            <br>
            <button type="submit">Login</button>
        </form>
    {{ end }}

    {{ with .Expenses }}
        <h2>Expenses</h2>
        <a href="{{ addQueryParam $.CurrentURL "from" (timestamp (addMonths $.To -2)) }}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-left-circle-fill"></i>
        </a>
        {{ formatTime (lastMonth $.To) "2006-01-02" }}
        <a href="{{ addQueryParam $.CurrentURL "from" $.Next }}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-right-circle-fill"></i>
        </a>

        {{ range .Currencies }}
            {{ $last_index := decrease (len .Intervals) }}
            <table class="table table-sm table-hover">
                <thead>
                    <tr>
                        <th>
                            <a href="{{ addQueryParam $.CurrentURL "currency" .CurrencyName }}"
                                    class="btn btn-light {{if eq (index $.Query "currency") .CurrencyName}}active{{end}}" tabindex="-1" role="button">
                                <i class="bi-arrow-through-heart-fill"></i>
                            </a>
                            Currency: {{ .CurrencyName }}
                        </th>
                        <th>{{ formatTime (index .Intervals $last_index) "2006-01-02" }}</th>
                    </tr>
                </thead>
                <tbody class="table-group-divider">
                    {{ $intervals := .Intervals }}
                    {{ range .Accounts }}
                    <tr>
                        {{ $accountID := .AccountID }}
                        <td>{{ .AccountName }}</td>
                        {{ $a := index .Amounts $last_index }}
                        <td {{ if and (eq $accountID "") (ne $a 0.0) }}class="table-danger"{{end}}>
                            {{ if eq $a 0.0 }}
                                0
                            {{ else }}
                                <a class="link-dark link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
                                    href="/web/transactions?from={{timestamp (index $intervals $last_index)}}&accountID={{$accountID}}">
                                        {{ money $a }}
                                </a>                    
                            {{ end }}
                        </td>
                        {{ end }}
                    </tr>
                </tbody>
                <tfoot class="table-secondary">
                    <td>
                        <strong>Total</strong>
                    </td>
                    <td>{{ money (index .Total $last_index) }}</td>
                </tfoot>                
            </table>
        {{ end }}
    {{ end }}

</main>

{{ template "footer.tpl" . }}