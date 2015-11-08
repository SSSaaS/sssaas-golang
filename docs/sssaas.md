# SSSaaS in Go
## String GetSecret(String[] API_Servers, String[] API_Tokens, String[] shares, timeout=300)
    API_Servers - List of API endpoints
    API_Tokens  - API Keys for matching endpoints
    share       - The known secret shares
    timeout     - Timeout to wait for tokens

    Combines the supplied secrets with the remote secrets using Shamir's
    Secret Sharing Algorithm. Blocks until all return and are combined.

    Returns the original secret.
