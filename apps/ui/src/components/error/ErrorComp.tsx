function ErrorComp({
  statusCode,
  errMsg,
}: {
  statusCode: number;
  errMsg: string;
}) {
  return (
    <div className="p-40 text-xl error-page">
      <div className="flex space-x-2">
        <span className="font-bold">{statusCode}</span>
        <span>|</span>
        <span>{errMsg}</span>
      </div>
    </div>
  );
}

export function Unauthorized() {
  return <ErrorComp statusCode={401} errMsg={"Unauthorized"} />;
}

export function Forbidden() {
  return <ErrorComp statusCode={403} errMsg={"Forbidden"} />;
}
