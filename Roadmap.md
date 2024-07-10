Have a look at:
https://github.com/gonzalezreal/swift-markdown-ui


Adjust max tokens anthropic.

Review types.

PDF review.

Read a file


Support multiple purposes.

Here is a revised version of your text for clarity:

---

The next stage of the project involves the following tasks:

1. **File Reading and Proofing:**
   - Enable the reading of a file. This file can be specified either as a default file or via a flag on the command line.
   - Send the contents of this file to the AI for proofing.

2. **Output Options:**
   - Display the proofed content on the screen.
   - Alternatively, output the proofed content to another file specified by a flag and the name of the file.

3. **Handling Additional Information:**
   - If the proofing prompt includes the field `request_additional_info` set to `true`, collect additional information.
   - This additional information can be gathered either by querying the user from the command line or by receiving it from another command line parameter.

---

This should make it clear to the AI what actions to take for each task.


Considerations:
Are their specific ways of passing in files to the API.
Different types of files.
May be handy to support the different roles so we can be clear.


Summarisation of PDF's will be done separately using the Assistants API.
