# Clear

**You can view the [Clear (Interpreted)](/clear/README.md) Folder to see the current stable implementation of Clear**

**You can also see examples of the language "talking" in the `/examples/` folder ([link1](./examples/00/00.log.md), [link2](./examples/01/01.log.md)). Addtionally, you can see the AST the program generated in the same folder (in JSON format)**

This repo contains tinkering on interpreters and compilers in Go

If I care about one of the subfolders that exists in this repo, than hopefully it should have a corresponding `README.md` file to explain what's going on there

I've been using multiple resources such as:
- [Writing an Interpreter in Go](https://interpreterbook.com/) by `Thorston Ball`
- [Writing a Compiler in Go](https://compilerbook.com/) also by `Thorston Ball`
- `Tyler Laceby`'s series on parsing with the Pikachu on the cover [link](https://www.youtube.com/watch?v=V77J9l8N-P8&list=PL_2VhOvlMk4XDeq2eOOSDQMrbZj9zIU_b)
- `Pixeled`'s series on writing a compiler in C++ [link](https://www.youtube.com/watch?v=vcSijrRsrY0&list=PLUDlas_Zy_qC7c5tCgTMYq2idyyT241qs)
- The classic [Crafting Interpreters](https://craftinginterpreters.com/)
- My project [GoDocs](https://github.com/ajtroup1/AdamTroup-430-Project), which operates similarly to Doxygen or Rust's documentation generator
  - Maybe the most I've learned during a project
  - Helped me realize that Parsers can be simple organizers that arrange information in a hierarchal, object-oriented manner
  - AST visualizers also help illustrate this point

Forever Clear