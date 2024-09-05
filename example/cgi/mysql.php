<h1>SQLRun Test</h1><hr>
<p>Show tables in sdnctl_nft</p><br>

<?php
// Retrieve the query parameter from the URL
$query = isset($_GET['query']) ? $_GET['query'] : '';

// Default to "show user()" if query is empty
if (empty($query)) {
    $query = 'select user()';
}

// Check the first word of the query
$firstWord = strtolower(strtok($query, ' '));
if ($firstWord !== 'select' && $firstWord !== 'show' && $firstWord !== 'describe') {
    echo 'Sorry read only actions';
    $query = 'select user()';
}

// If the first word is "select", add a limit to the query
if ($firstWord === 'select') {
    if (stripos($query, 'limit') === false) {
        // No limit specified, add limit 100
        $query .= ' LIMIT 100';
    } else {
        // Limit specified, ensure it does not exceed 100
        $query = preg_replace_callback('/limit\s+(\d+)/i', function ($matches) {
            return 'LIMIT ' . min(100, (int)$matches[1]);
        }, $query);
    }
}

// Command to execute
$command = sprintf('sqlrun -c /etc/sqlrun/sdncp-prod.conf -e "%s" -o html', $query);
echo "<p><b>Command:</b> $command</p>";

// Execute the command
exec($command, $output, $return_var);

// Check if the command was successful
if ($return_var !== 0) {
    echo 'Error executing command';
} else {
    // Display the output
    echo implode("\n", $output);
}
?>
<p><b>Note:</b>Queries are limited to 100 rows.</b>
</body></html>