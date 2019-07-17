# Contributing

All contributions are welcome, just follow our contribution guide :)

If you need a new feature in our code base, simply write your feature and send a merge request. It's encouraged to talk with the system owner beforehand.

The whole codebase is dockerized so you should be able to run the codebase easily simply using docker. For more information about how to run please read `README.md`

# Development
1. All the (code) parts that you've changed or added must have tests.
2. Test coverage should not decrease by merging your code.
3. Commit messages should be self explanatory. It should answer following questions:
    * Why is this change necessary?
    * How does it address the issue?
    * What side effects does this change have?
4. None of your (final) commits should break CI tests.
    * There are some minor exceptions but believe me, there's a high probability that your case is not one of them!
5. Your (final) branch should have just the right number of commits not too many, not too few.
    * Logically relevant changes SHOULD get committed together.
    * Logically irrelevant changes SHOULD NOT get committed together.

# Merge Request Process
1. Rebase your code to current master's head (Do not merge master into your branch, this ruins your branch's history).
2. Push your branch upstream and create a merge request.
3. Keep in mind that all merge request must satisfy following criteria to get merged:
    * You should test your code vigorously. It's assumed that your merge request does not break anything and works the way it's supposed to.
    * Merging your merge request should not decrease test coverage.
    * If you're making a change that needs a new side service (like redis) or changes `config.template.json` you SHOULD mention it in your merge request.
4. Please resolve merge request's issues in less than 7 working days. Keep in mind that inactive merge requests (inactive for more than 14 days) will be closed.
