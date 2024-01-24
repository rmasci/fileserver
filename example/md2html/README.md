# Verbose

Allows you to easily add in statements to your code that allow you to trace what your code is doing or what a variable,
is equal to at a certain point.  This also prints the line number so you can quickly find what file / line we're on. Good
for when you've not looked at the code in a year or more. 

For example, lets say in your code you're doing a database query, no error but it gives you unexpected results. Use 
verbose to print the query, maybe the query isn't right in this instance. 

The way it is intended is you can leave these in the code. Use it in conjunction with flag or pflag so that you can just
invoke it. Normally the user will not see the verbose statements, but a -v or --verbose when running, and they 
show up.

This will also print to a file, or to stderr if you want it to.  This allows  you to have your end user run it in verbose
mode, log it to a file that they can send you in an email allowing you to evaluate what's going on with your code.

Example Code:
```cgo
func main() {
	i := 0
	var testme string
	var help bool
	// Notice the formatting is the same as if you used the Linux date command. 
	verb := verbose.New("%A %B %Y, %I:%M:%S %P %Z")
	var flagset pflag.FlagSet
	flagset.BoolVarP(&verb.V, "verbose", "v", false, "Verbose Mode")
	flagset.BoolVarP(&help, "help", "h", false, "Help")
	flagset.StringVarP(&testme, "string", "s", "", "String to print out.")
	flagset.MarkHidden("verbose")
	flagset.Parse(os.Args[1:])
	
	verb.Printf("%s started\n",os.Args[0])
	verb.Printf("Remaining Args passed: %v\n",flagset.Args())
	... lots of code ...
	// verb.
	verb.Println("Query Database for user status:")
	verb.Println("Query:",query)
	... lots of code ...
	verb.Println("Rest Call to www.somewhere.com")
```
When the user runs a program like this if the -v or --verbose is not passed as an option, the verb.Print statements will not print.

The default output is to os.Stdout - but this can be changed. If in the above example we wanted to write the output to a file. or to stderr we can do that two ways.
1. set verb.Out:=*os.Writer
```cgo
// set message to file
w:=os.Create("verbose.txt")
verb.Out=w
verb.Println("This line goes to a file")
```
2. Use verb.Fprintf(os.Writer, "Your message goes here.")
```cgo
// Set message to stderr:
verb.Fprintf(os.Stderr, "My Message Goes Here")
```


