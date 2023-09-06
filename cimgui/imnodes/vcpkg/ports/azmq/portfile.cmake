vcpkg_from_github(
    OUT_SOURCE_PATH SOURCE_PATH
    REPO zeromq/azmq
    REF 6bb101eecb357ad9735ebc36e276b7526652d42d # commit on 2019-05-01 
    SHA512 18812fd73c3c57aca8b17c2df0af01afb7864253d6ac6ce042e01ef04e81dd454438c9d01f5bb0fd7a4a445082a2eb3dd334ca1e05dafbe4ba602cad95ce7134
    HEAD_REF master
)

file(COPY ${SOURCE_PATH}/azmq DESTINATION ${CURRENT_PACKAGES_DIR}/include/)

file(INSTALL
    ${SOURCE_PATH}/LICENSE-BOOST_1_0
    DESTINATION ${CURRENT_PACKAGES_DIR}/share/${PORT} RENAME copyright)