
# go-tip
This library is inspired by the [TiP Scala library](https://github.com/guardian/tip) used by The Guardian. This should be treated as a proof of concept at this stage and has not been reviewed or performance-optimised.

The basic principle is to use this library to track which of your applications critical paths have been successfully traversed following a merge and deploy of your code. Once all paths have been successfully reached, the most-recent pull request is marked as Verified, giving developers confidence that they have not broken anything fundamental.

### Usage

Import the packahge:

	import "github.com/martinroddam/go-tip"

Add a `paths.yaml` file into the root of your project. See `/examples`.

Add a `secret.yaml` file into the `secrets` directory located at the root of your project. This file should detail the github owner e.g. `utilitywarehouse`, the repo name, and an Access token which has at least `public_repo` scope. **Remember to encrypt!**. See `/examples`.

Add verification points throughout your app as follows:

    gotip.Verify("Order Submitted")

Note: The string value provided is a label for the particular critical path in your application that you have successfully reached. The labels may contain spaces and are case-sensitive. Labels can be given a colour to make them stand out more - see Github label settings.

### TODO

Allow users to pass in an environment parameter so that the label can say `Verified in DEV`, or `Verified in PROD` etc
Consider checking that the deployed sha matches the Pull Request sha before Verifying paths.