import { useTranslation } from "react-i18next";
import Head from "next/head";

import { SignUpPage } from "@app/views/admin/SignUpPage";

export default function LoginPage() {
  const { t } = useTranslation();

  return (
    <div className="flex min-h-screen max-h-screen overflow-y-auto flex-col justify-center bg-gradient-to-tr from-mineshaft-600 via-mineshaft-800 to-bunker-700 px-6">
      <Head>
        <title>{t("common.head-title", { title: t("signup.title") })}</title>
        <link rel="icon" href="/gsoc2.ico" />
        <meta property="og:image" content="/images/message.png" />
        <meta property="og:title" content={t("signup.og-title") ?? ""} />
        <meta name="og:description" content={t("signup.og-description") ?? ""} />
      </Head>
      <SignUpPage />
    </div>
  );
}
