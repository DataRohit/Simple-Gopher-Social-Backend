{{- /* gotype: gopher-social-backend-server/pkg/mailer/sendOAuthWelcomeMail.OAuthWelcomeEmailData */ -}}

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Our Service</title>
</head>
<body>
    <h1>Welcome to Our Service!</h1>
    <p>Hi {{.Email}},</p>
    <p>Thank you for registering with your {{.Provider}} account. We're excited to have you on board!</p>
    <p>Thank you!</p>
</body>
</html>
