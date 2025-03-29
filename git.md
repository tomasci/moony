Private projects

1. switch to private/main branch (private is private repo remote)
2. work there
3. when update is ready, switch to origin/main 
4. commit "public" part in main
5. push origin/main
6. switch back to private/main
7. merge origin/main to private/main
8. commit "private" part in private/main
9. push private/main

Pull requests and other branches can be used, if you don't want
to commit to main branch

-- 

This if mostly for me to keep private projects and everything else 
out of public version of Moony server 