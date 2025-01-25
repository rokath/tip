//! @file tipTable.c
//! @brief Generated code - do not edit!

#include "tipTable.h"

tip_t TipTable[] = {
    //  by, sz, pattern                                            // occurances
    { 0x52, 8, {0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xaa} }, //         13
    { 0x12, 4, {0xff, 0xff, 0xee, 0xcc} },                         //        876
    { 0x56, 4, {0xfa, 0xaf, 0x00, 0xaa} },                         //        123
    { 0x04, 2, {0xfe, 0xff} },                                     //      12345
    { 0x34, 2, {0xff, 0xff} },                                     //       1234
    { 0x56, 2, {0xfa, 0xaf} }                                      //        123
};

int const TipTableLength = 6; //!< Number of TipTableEntries, cannot exceed 127.
