# Problems with the Source Code (Ranked)

1. Secret uploaded in the source code `SECRET_KEY`
2. Secret uploaded in `schema.sql` — contains an `INSERT` showing the admin password.
3. SQL injection is possible because user input is inserted directly into SQL statements, for example:
   - `login, register, api/search`
   - Parameterized queries are more secure.
4. Passwords are hashed with MD5 and without salt. Use stronger hashing instead and combine it with salt.
5. Empty test `search()`
6. When sending a request to `/api/login` with a correct username but an incorrect password, the server responds that the password is invalid. This leaks to the client that a user with the given username exists, which can make it easier to gain access to other users' accounts.
7. When sending a request to `/api/register` with a unique username but an email that is not unique, the server returns an internal server error 500, which does not follow the OpenAPI specification.
8. The variable `one` in `query_db()` is confusing. Its usage indicates that a more helpful name could be `return_only_first_row`
9. In addition, *cur* could possibly be called `response`, and `rv` could be called `row`, in `query_db()`.
- `q` in `search()` is also a confusing variable name. Does it mean (SQL) query or a question from the user?
- `get_user_id()` should be named something like `get_user_id_for_username()`.
10. Tests are generally weak.
11. Dependency pinning in `requirements.txt` is outdated.
12. A `</li>` tag is missing in `layout.html`.
13. In the `schema.sql` file, there is a comment on line 18 explaining enums. Is that really strictly necessary?
14. The footer floats in the middle of the page instead of being a footer.
