"use client";

import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import { Input, Stack, Text } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import { useState } from "react";
import { useForm } from "react-hook-form";
import z from "zod";

const schema = z.object({
  name: z.string().min(1),
  key: z.string().min(1),
});

type FormData = z.infer<typeof schema>;

export default function Form({ address }: { address: string }) {
  const [success, setSuccess] = useState(false);
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({ resolver: zodResolver(schema) });

  const onSubmit = async ({ name, key }: FormData) => {
    const validate = schema.safeParse({ name, key });
    if (validate.success) {
      // example of api call from client
      const res = await axios
        .post(`/api/send`, {
          address,
          name,
          key,
        })
        .catch((e) => {
          if (axios.isAxiosError(e)) {
            console.log("HTTPS req is not working");
          }
        });
      setSuccess(true);
      console.log(res);
    }
  };
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Stack
        gap="8"
        width="100%"
        height={"fit-content"}
        css={{ "--field-label-width": "96px" }}
      >
        <Field
          orientation="vertical"
          label="Name"
          invalid={Boolean(errors.name)}
          errorText={errors.name?.message}
        >
          <Input placeholder="John" flex="1" p={3} {...register("name")} />
        </Field>
        <Field
          orientation="vertical"
          label="Key"
          invalid={Boolean(errors.key)}
          errorText={errors.key?.message}
        >
          <Input placeholder="Doe" flex="1" p={3} {...register("key")} />
        </Field>
        <Button type="submit" loading={isSubmitting}>
          Submit
        </Button>
        {success && (
          <Text fontSize={"md"} color={"WindowText"}>
            Root:
            <br />
            0xc5b7f90e2a9d7622cfab1f8bb3becd1317bd4af85097afdc6815c03d15388e61{" "}
            <br /> New Verifier deployed with correct root <br /> <br />
            Proof: [
            <br />
            &apos;0xc9f81d534037cca28de7d2aa8c62e5d6b75d5b58ccc4a265138671179cc4d447&apos;,
            <br />
            <br />
            &apos;0x55d8939601c57fb83ef75669c81173aa15197063dc46fa4b653b89862dcde003&apos;,
            <br />
            ]
            <br />
            <br />
            Attempting to store API key... <br />
            Transaction sent: <br /> <br />
            Transaction successful!
          </Text>
        )}
      </Stack>
    </form>
  );
}
