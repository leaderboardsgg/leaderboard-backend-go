# This is a sample GraphQL Server!
## How do I run it?
`make gql_run`

## How do I see the output?
- You can go to `localhost:3030/graphiql`.
- You can issue HTTP Requests to `localhost:3030/graphql/http`.

## How do I use graphiql?
You enter the query you want on the left side, and click play!

## Do you have any example queries?
```
{
  games(titleRegex: "Great"){
    title
    runs {
      time
      runner {
        name
      }
    }
  }
}
```
This nesting shows the power and flexibility of GraphQL, I think.

# Call to action!
## What can you do?
- Make a real server with real types.
- Delete this sample.
