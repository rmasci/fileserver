<?php
// Check if form is submitted
if ($_SERVER["REQUEST_METHOD"] == "POST") {
    // Retrieve form data
    $num1 = isset($_POST['num1']) ? (float)$_POST['num1'] : 0;
    $num2 = isset($_POST['num2']) ? (float)$_POST['num2'] : 0;
    $operation = isset($_POST['operation']) ? $_POST['operation'] : 'add';

    // Perform the selected operation
    switch ($operation) {
        case 'add':
            $result = $num1 + $num2;
            break;
        case 'subtract':
            $result = $num1 - $num2;
            break;
        case 'multiply':
            $result = $num1 * $num2;
            break;
        case 'divide':
            if ($num2 != 0) {
                $result = $num1 / $num2;
            } else {
                $result = 'Error: Division by zero';
            }
            break;
        default:
            $result = 'Invalid operation';
            break;
    }
} else {
    $num1 = $num2 = $result = $operation = '';
}
?>

<!DOCTYPE html>
<html>
<head>
    <title>PHP Calculator</title>
</head>
<body>
    <h1>PHP Calculator</h1>
    <form method="post" action="">
        <label for="num1">Number 1:</label>
        <input type="text" id="num1" name="num1" value="<?php echo htmlspecialchars($num1); ?>"><br><br>
        <label for="num2">Number 2:</label>
        <input type="text" id="num2" name="num2" value="<?php echo htmlspecialchars($num2); ?>"><br><br>
        <label for="operation">Operation:</label>
        <select id="operation" name="operation">
            <option value="add" <?php if ($operation == 'add') echo 'selected'; ?>>Add</option>
            <option value="subtract" <?php if ($operation == 'subtract') echo 'selected'; ?>>Subtract</option>
            <option value="multiply" <?php if ($operation == 'multiply') echo 'selected'; ?>>Multiply</option>
            <option value="divide" <?php if ($operation == 'divide') echo 'selected'; ?>>Divide</option>
        </select><br><br>
        <input type="submit" value="Calculate">
    </form>

    <?php if ($_SERVER["REQUEST_METHOD"] == "POST"): ?>
        <h2>Result: <?php echo htmlspecialchars($result); ?></h2>
    <?php endif; ?>
</body>
</html>