using Aifunc;
using Aifunc.Summarize;

// var config = new AIFuncConfig
// {
//     BaseUrl = "https://your-api-endpoint/v1",
//     Model = "your-model-name",
//     ApiKey = "your-api-key",
//     MaxRetries = 3,
// };

// To use a real model, replace the line below with the commented config above.
var config = new AIFuncConfig { Mock = true };

if (config.Mock)
{
    Console.WriteLine("Notice: You are using mock mode for offline testing. " +
        "Configure a real model for the full experience. Continuing with mock responses...");
}

var text =
    "The James Webb Space Telescope captured its first full-color images in July 2022, " +
    "revealing thousands of galaxies in a patch of sky smaller than a grain of sand held " +
    "at arm's length. The images show galaxies as they appeared over 13 billion years ago, " +
    "providing a glimpse into the early universe shortly after the Big Bang.";

var result = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput(text, 30));

Console.WriteLine("Original  : " + text);
Console.WriteLine("Summary   : " + result.Summary);
Console.WriteLine("Word count: " + result.WordCount);
