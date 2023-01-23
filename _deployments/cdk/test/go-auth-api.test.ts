import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import { GoAuthApiStack as TestStack } from '../lib/go-auth-api-stack';
import { context } from '../cdk.json';

// example test. To run these tests, uncomment this file along with the
// example resource in lib/go-auth-api-stack.ts
test('snapshot test', () => {
    const env = "stage"
    const app = new cdk.App({ context: { env, [env]: context } });
    // WHEN
    const stack = new TestStack(app, 'TestStack');
    // THEN
    const template = Template.fromStack(stack);

    expect(template).toMatchSnapshot();
});
