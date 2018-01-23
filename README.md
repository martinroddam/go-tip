
# go-tip
This library is inspired by the [TiP Scala library](https://github.com/guardian/tip) used by The Guardian.

The basic principle is to use this library to track which of your applications critical paths have been successfully traversed following a merge and deploy of your code. Once all paths have been successfully reached, the most-recent pull request is marked as Verified, giving developers confidence that they have not broken anything fundamental.

### Usage

Import the packahge:

	import "github.com/martinroddam/go-tip"

Add a `paths.yaml` file into the root of your project. See `/examples`.

Add verification points throughout your app as follows:

    gotip.Verify("Order Submitted")

Note: The string value provided is a label for the particular critical path in your application that you have successfully reached. The labels may contain spaces and are case-sensitive. 