{{- /* gotype: gopher-social-backend-server/pkg/mailer/accountActivationMail.ActivationEmailData */ -}}

<!DOCTYPE html>
<html>
<head>
    <title>Activate Your Account</title>
</head>
<body>
    <h1>Activate Your Account</h1>
    <p>Hi {{.Email}},</p>
    <p>Please activate your account by visiting the following link:</p>
    <a href="http://localhost:8080/auth/activate/{{.Token}}">Activate Account</a>
    <p>The link expires at {{.Expiration}}.</p>
    <p>Thank you!</p>
</body>
</html>
