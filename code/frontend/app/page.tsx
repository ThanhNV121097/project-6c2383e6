import PersistEditableGreeting from "../components/PersistEditableGreeting";
import { readGreetingFromServer } from "../lib/persist-editable-greeting";

export default async function Page() {
  const greeting = await readGreetingFromServer();

  return (
    <main>
      <PersistEditableGreeting initialGreeting={greeting} />
    </main>
  );
}
