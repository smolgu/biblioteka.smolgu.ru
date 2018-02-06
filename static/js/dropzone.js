var myDropzone = new Dropzone(".dropzone", {  });
myDropzone.on("success", function( file, result ) {
  // the file parameter is https://developer.mozilla.org/en-US/docs/DOM/File
  // the result parameter is the result from the server

  // [success code here]
  console.log(file);
  console.log(result);
});