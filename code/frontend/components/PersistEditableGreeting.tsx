"use client";

import { FormEvent, useEffect, useRef, useState } from "react";
import {
  mockGreeting,
  readStoredGreeting,
  saveStoredGreeting,
  type Greeting,
} from "../lib/mock/persist-editable-greeting";
import styles from "./PersistEditableGreeting.module.css";

type PersistEditableGreetingProps = {
  initialGreeting?: Greeting;
};

export default function PersistEditableGreeting({
  initialGreeting = mockGreeting,
}: PersistEditableGreetingProps) {
  const [greeting, setGreeting] = useState(initialGreeting.text);
  const [inputValue, setInputValue] = useState(initialGreeting.text);
  const [status, setStatus] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const storedGreeting = readStoredGreeting();
    setGreeting(storedGreeting.text);
    setInputValue(storedGreeting.text);
  }, []);

  function saveGreeting(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextGreeting = inputValue.trim();

    if (!nextGreeting) {
      setStatus("Enter a greeting.");
      inputRef.current?.focus();
      return;
    }

    const savedGreeting = saveStoredGreeting(nextGreeting);
    setGreeting(savedGreeting.text);
    setInputValue(savedGreeting.text);
    setStatus("Saved.");
  }

  return (
    <section className={styles.section} aria-labelledby="greeting">
      <h1 id="greeting" className={styles.heading}>
        {greeting}
      </h1>
      <form className={styles.form} onSubmit={saveGreeting} noValidate>
        <label className={styles.label} htmlFor="greeting-input">
          Greeting
        </label>
        <input
          ref={inputRef}
          id="greeting-input"
          name="greeting"
          type="text"
          value={inputValue}
          autoComplete="off"
          required
          className={styles.input}
          onChange={(event) => setInputValue(event.target.value)}
        />
        <button className={styles.button} type="submit">
          Save
        </button>
      </form>
      <p className={styles.status} id="status" role="status" aria-live="polite">
        {status}
      </p>
    </section>
  );
}
