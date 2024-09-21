{{- /* gotype: gopher-social-backend-server/pkg/mailer/passwordResetMail.PasswordResetEmailData */ -}}

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
</head>
<body>
    <h1>Password Reset Request</h1>
    <p>Hi {{.Email}},</p>
    <p>You requested to reset your password. Please click the link below to reset your password:</p>
    <a href="http://localhost:8080/auth/reset-password/{{.Token}}">Reset Password</a>
    <p>This link will expire in {{.Expiration}}.</p>
    <p>If you didn't request this, please ignore this email.</p>
    <p>Thank you!</p>
</body>
</html>
