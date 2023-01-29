import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import { GoAuthApiStack as TestStack } from '../lib/go-auth-api-stack';
import { context } from '../cdk.json';

// example test. To run these tests, uncomment this file along with the
// example resource in lib/go-auth-api-stack.ts
test('snapshot test', () => {
    const env = "stage"
    const app = new cdk.App({ context: { env, ...context } });
    // WHEN
    const stack = new TestStack(app, 'TestStack');
    // THEN
    const template = function () {
        const template = Template.fromStack(stack).toJSON();
        // IGNORE
        for (const resourceKey in template["Resources"]) {
            const resource: { [key: string]: any } = template["Resources"][resourceKey]
            const resourceType: string = resource["Type"]
            const properties = resource["Properties"]
            if (resourceType == "AWS::Route53::RecordSet" && properties["Type"] == "A") {
                resource["Properties"]["Name"] = ""
            }
        }
        return Template.fromJSON(template)
    }()

    expect(template).toMatchSnapshot();
});
